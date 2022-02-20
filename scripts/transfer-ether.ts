import { ethers } from 'ethers';
import * as dotenv from 'dotenv';

(async function () {
  dotenv.config();
  const provider = ethers.getDefaultProvider('http://localhost:8545');

  const network = await provider.getNetwork();

  console.log('Network name=', network.name);
  console.log('Network chain id=', network.chainId);

  const privateKey = process.env.PRIVATE_KEY;

  if (!privateKey) {
    throw new Error('PRIVATE_KEY environment variable is undefined');
  }

  const signer = new ethers.Wallet(privateKey, provider);

  signer.sendTransaction({
    to: '0x2bb6335AC37c468c626D18C9915A8Cc7c36D76e7',
    value: ethers.utils.parseUnits('100', 'ether'),
  });
})();
