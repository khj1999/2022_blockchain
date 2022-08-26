'use strict';

// 1. server set
var express = require('express')
var path = require('path')
var fs = require('fs')
var bodyParser = require('body-parser');


// connection.json °´Ã¼È­
const ccpPath = path.resolve(__dirname, "..", "..", "..", "fabric-samples", "test-network","organizations","peerOrganizations", "org1.example.com", "connection-org1.json");
const ccp = JSON.parse(fs.readFileSync(ccpPath,"utf8"));

// 2. fabric connect set
const { Gateway, Wallets } = require('fabric-network');
const FabricCAServices = require('fabric-ca-client');

const { buildCAClient, registerAndEnrollUser, enrollAdmin } = require('./CAUtil.js');
const { buildCCPOrg1, buildWallet } = require('./AppUtil.js');

const mspOrg1 = 'Org1MSP';
const walletPath = path.join(__dirname, 'wallet')
const caClient = buildCAClient(FabricCAServices, ccp, 'ca.org1.example.com')


// 3. midlle ware set
var app = express();

app.use('/UI',express.static(path.join(__dirname,'UI')));
app.use(express.static(path.join(__dirname,"UI")));
app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());


app.post("/trade",async(req, res) =>{
    const product = req.body.product
    const trade_id = req.body.trade_id;
    const seller_name = req.body.seller_name;

    console.log("/trade post start -- ", product, trade_id, seller_name);
    const gateway = new Gateway();

    try{
        const wallet = await buildWallet(Wallets, walletPath);
        await gateway.connect(ccp, {
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true }
        });
        const network = await gateway.getNetwork("mychannel");
        const contract = network.getContract("key");
        await contract.submitTransaction('Register', trade_id, product, seller_name);

    } catch(error){
        var result = `{"result":"fail", "message":"tx has NOT submitted"}`;
        var obj = JSON.parse(result);
        console.log("/trade end -- failed ", error);
        res.status(200).send(obj);
        return;
    } finally{
        gateway.disconnect();
    }
    var result = `{"result":"success", "message":"tx has submitted"}`;
    var obj = JSON.parse(result);
    console.log("/asset end -- success");
    res.status(200).send(obj);
});


app.get("/trade", async(req,res)=>{
    const trade_id = req.query.trade_id;
    console.log("/trade get start -- ", trade_id);

    const gateway = new Gateway();

    try{
        const wallet = await buildWallet(Wallets, walletPath);

        await gateway.connect(ccp,{
            wallet,
            identity: 'appUser',
            discovery: { enabled: true, asLocalhost: true}
        });
        const network = await gateway.getNetwork("mychannel");
        const contract = network.getContract("key");
        var result = await contract.evaluateTransaction('GetTrade',trade_id);
        var result = `{"result":"success", "message":${result}}`;
        console.log("/trade get end -- success", result);
        var obj = JSON.parse(result);
        res.status(200).send(obj);
    } catch (error){
        var result = `{"result":"fail","message":"Get has a error"}`;
        var obj = JSON.parse(result);
        console.log("/trade get end -- faild ", error);
        res.status(200).send(obj);
        return;
    } finally{
        gateway.disconnect();
    }

});

app.post("/trade/tx",async(req, res) =>{
    const trade_id = req.body.trade_id
    const product = req.body.product;
    const address = req.body.address;
    const purchaser_name = req.body.purchaser_name;

    console.log("/trade application start -- ", trade_id, product, address, purchaser_name);
    const gateway = new Gateway();

    try{
        const wallet = await buildWallet(Wallets, walletPath);
        await gateway.connect(ccp,{
            wallet,
            identity: "appUser",
            discovery: { enabled: true, asLocalhost: true }
        });
        const network = await gateway.getNetwork("mychannel");
        const contract = network.getContract("key");
        await contract.submitTransaction('Request', trade_id, product, address, purchaser_name);

    } catch(error){
        var result = `{"result":"fail", "message":"tx has NOT submitted"}`;
        var obj = JSON.parse(result);
        console.log("/application end -- failed ", error);
        res.status(200).send(obj);
        return;
    } finally{
        gateway.disconnect();
    }
    var result = `{"result":"success", "message":"${product} application , tx has submitted at ${purchaser_name}."}`;
    var obj = JSON.parse(result);
    console.log("/application end -- success");
    res.status(200).send(obj);
});

app.listen(3000, ()=>{
    console.log('Express server is started: 3000');
});
