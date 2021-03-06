import appRootPath from 'app-root-path';
import * as child_process from 'child_process';
import * as dotenv from 'dotenv';
import { ContractFactory, ContractInterface, ethers } from 'ethers';
import CryptoKoi from '../artifacts/contracts/CryptoKoi.sol/CryptoKoi.json';

(async function init() {
  dotenv.config();

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

  const deployContract = (
    abi: ContractInterface,
    bytecode: string,
    signer: ethers.Signer,
    options: {
      name: string;
      symbol: string;
      baseURI: string;
      priceInMatic: number;
      gasPrice: number;
    },
  ): Promise<ethers.Contract> => {
    const contract = new ContractFactory(abi, bytecode, signer);
    return contract.deploy(
      options.name,
      options.symbol,
      options.baseURI,
      ethers.utils.parseEther(options.priceInMatic.toString()),
      {
        maxFeePerGas: options.gasPrice,
        maxPriorityFeePerGas: options.gasPrice,
      },
    );
  };

  const privateKey = process.env.PRIVATE_KEY;

  if (!privateKey) {
    throw new Error('PRIVATE_KEY environment variable is undefined');
  }

  const gasPrice = ethers.utils.parseUnits('250', 'gwei').toNumber();
  const provider = ethers.getDefaultProvider(url);
  const signer = new ethers.Wallet(privateKey, provider);
  const c = new ContractFactory(
    CryptoKoi.abi,
    CryptoKoi.bytecode,
    signer,
  );

  const deploymentData = c.interface.encodeDeploy([
    'CryptoKoi',
    'CK',
    'https://api.crypto-koi.io/v1/tokens/',
    // MATIC calculation
    ethers.utils.parseEther((1.66).toString()),
  ]);

  const gasUnitsNeeded = (
    await provider.estimateGas({
      data: deploymentData,
    })
  ).toNumber();

  const pricePerUnitInGwei = gasPrice;

  const gasEstimateCoin =
    (gasUnitsNeeded * pricePerUnitInGwei) / (1000000000 * 1000000000);

  console.log(`Gas estimate: ${gasEstimateCoin.toString()} MATIC


Gas Units needed: ${gasUnitsNeeded}
Price per unit in Wei: ${pricePerUnitInGwei}

${(1000000000 * 1000000000).toString()} Gwei = 1 Matic
`);

  const contract = await deployContract(
    CryptoKoi.abi,
    CryptoKoi.bytecode,
    signer,
    {
      name: 'CryptoKoi',
      symbol: 'CK',
      baseURI: 'https://api.crypto-koi.io/v1/tokens/',
      // MATIC calculation
      priceInMatic: 1.6, // around 1,99???
      gasPrice: gasPrice,
    },
  );

  console.log('DEPLOY TRANSACTION', contract.deployTransaction);

  const contractAddress = (await contract.deployTransaction.wait())
    .contractAddress;

  console.log(
    'Save the contract address as environment variable. Copy the following line and paste it into the .env file. The client application will need it as well',
  );
  console.log('CONTRACT_ADDRESS=' + contractAddress);
})();
