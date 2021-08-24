/**
 * determine if the current build is for the TestNet or for Production
 * @returns {bool} if the current build is for the testnet
 */
export function isTestNet() {
  if (process.env.REACT_APP_GDS_IS_TESTNET) {
    return Boolean(JSON.parse(process.env.REACT_APP_GDS_IS_TESTNET));
  }
  return false;
}