import * as dotenv from "dotenv";
import Hello from "../artifacts/contracts/Hello.sol/Hello.json";

import {join} from "path";
import Web3 from "web3";
import * as  child_process from "child_process";
import appRootPath from "app-root-path";
import Contract from "web3/eth/contract";
import Eth from "web3/eth";

dotenv.config({
    path: join(process.cwd(), "..", ".env"),
});

const deployContract = async (
    web3: Web3,
    abi: any | any[],
    data: string,
    from: string
) => {
    const deployment = new web3.eth.Contract(abi).deploy({ data });
    const gas = await deployment.estimateGas();
    const {
        options: { address: contractAddress },
    } = await deployment.send({ from, gas });
    return new web3.eth.Contract(abi, contractAddress);
};


const opts: child_process.ExecSyncOptions = {
    cwd: appRootPath.toString(),
    stdio: "inherit"
}
// compile the smart contract.
child_process.execSync("npx hardhat compile", opts);

const url = process.env.CHAIN_URL;

if (!url) {
    throw new Error("CHAIN_URL environment variable is undefined.")
}

const privateKey = process.env.PRIVATE_KEY

if (!privateKey) {
    throw new Error("PRIVATE_KEY environment variable is undefined")
}

(async function fn() {
    const web3 = new Web3(new Web3.providers.HttpProvider(url));

    const {address} = web3.eth.accounts.privateKeyToAccount(privateKey)
    
    const contract = await deployContract(web3, Hello.abi, Hello.bytecode, address);
    console.log(contract)
    console.log(await contract.methods.sayHello("React").call())
})();


