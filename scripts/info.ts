import { ethers } from 'ethers';
import * as dotenv from 'dotenv';

(async function () {
  dotenv.config();
  const provider = ethers.getDefaultProvider(process.env.CHAIN_URL);

  const network = await provider.getNetwork();

  console.log('Network name:', network.name);
  console.log('Network chain id:', network.chainId);
})();
