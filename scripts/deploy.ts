import appRootPath from 'app-root-path';
import * as child_process from 'child_process';
import * as dotenv from 'dotenv';
import { join } from 'path';
import Web3 from 'web3';
import CryptoKoi from '../artifacts/contracts/CryptoKoi.sol/CryptoKoi.json';

dotenv.config({
  path: join(process.cwd(), '..', '.env'),
});

const deployContract = async (
  web3: Web3,
  abi: any | any[],
  data: string,
  from: string,
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
  stdio: 'inherit',
};
// compile the smart contract.
child_process.execSync('npx hardhat compile', opts);

const url = process.env.CHAIN_URL;

if (!url) {
  throw new Error('CHAIN_URL environment variable is undefined.');
}

const privateKey = process.env.PRIVATE_KEY;

if (!privateKey) {
  throw new Error('PRIVATE_KEY environment variable is undefined');
}

(async function fn() {
  const web3 = new Web3(new Web3.providers.HttpProvider(url));

  const { address } =
    web3.eth.accounts.privateKeyToAccount(privateKey);

  const contract = await deployContract(
    web3,
    CryptoKoi.abi,
    CryptoKoi.bytecode,
    address,
  );
})();
