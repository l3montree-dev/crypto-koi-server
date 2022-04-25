import appRootPath from 'app-root-path';
import * as child_process from 'child_process';
import * as dotenv from 'dotenv';
import {
  Contract,
  ContractFactory,
  ContractInterface,
  ethers,
} from 'ethers';
import CryptoKoi from '../artifacts/contracts/CryptoKoi.sol/CryptoKoi.json';

(async function init() {
  dotenv.config();

  const opts: child_process.ExecSyncOptions = {
    cwd: appRootPath.toString(),
    stdio: 'inherit',
  };

  const url = process.env.CHAIN_URL;

  if (!url) {
    throw new Error('CHAIN_URL environment variable is undefined.');
  }

  const contractAddress = process.env.CONTRACT_ADDRESS;
  if (!contractAddress) {
    throw new Error(
      'CONTRACT_ADDRESS environment variable is undefined.',
    );
  }

  const privateKey = process.env.PRIVATE_KEY;

  if (!privateKey) {
    throw new Error('PRIVATE_KEY environment variable is undefined');
  }

  const provider = ethers.getDefaultProvider(url);
  const balance = await provider.getBalance(contractAddress);
  const signer = new ethers.Wallet(privateKey, provider);
  const c = new Contract(contractAddress, CryptoKoi.abi, signer);

  console.log('BALANCE', ethers.utils.formatEther(balance), 'MATIC');

  const tx = await c.withdrawAll();
  console.log('TX', tx.hash);
  const receipt = await tx.wait();
  console.log('RECEIPT', receipt);
})();
