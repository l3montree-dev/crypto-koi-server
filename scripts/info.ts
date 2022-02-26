import { ethers } from 'ethers';
import * as dotenv from 'dotenv';

(async function () {
  dotenv.config();
  const provider = ethers.getDefaultProvider(
    'https://ropsten.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161',
  );

  const network = await provider.getNetwork();

  console.log('Network name:', network.name);
  console.log('Network chain id:', network.chainId);
})();
