import * as dotenv from 'dotenv';
import { Contract, ethers } from 'ethers';
import CryptoKoi from '../artifacts/contracts/CryptoKoi.sol/CryptoKoi.json';

(async function init() {
  dotenv.config();

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
  const signer = new ethers.Wallet(privateKey, provider);

  const contract = new Contract(
    contractAddress,
    CryptoKoi.abi,
    signer,
  );

  console.log(
    await contract.setBaseURI(
      'https://dev.api.crypto-koi.io/v1/tokens/',
    ),
  );
})();
