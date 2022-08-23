import { Text } from '@chakra-ui/react';
import { t } from '@lingui/macro';
import {
  getAppGitVersion,
  getAppVersionNumber,
  getBffAndGdsVersion,
  isProdEnv
} from 'application/config';
import { useEffect, useState } from 'react';

function Version() {
  const appVersion = getAppVersionNumber();
  const appGitVersion = getAppGitVersion();
  const [bffAndGdsVersion, setBffAndGdsVersion] = useState<any>();

  const fetchAsyncBffAndGdsVersion = async () => {
    const request = await getBffAndGdsVersion();
    if (request) {
      setBffAndGdsVersion(request.version);
    }
  };

  useEffect(() => {
    fetchAsyncBffAndGdsVersion();
  }, []);

  if (!isProdEnv) return null;

  //  log this out in the console
  console.log('appVersion', appVersion);
  console.log('gitRevision', appGitVersion);
  console.log('bffAndGdsVersion', bffAndGdsVersion);

  return (
    <Text width="100%" textAlign="center" color="white" fontSize="12" pt={1}>
      {appVersion && <Text as="span">{t`App version ${appVersion}`} - </Text>}
      {appGitVersion && (
        <Text as="span" data-testid="git-revision">
          {t`Git Revision`} {appGitVersion} -{' '}
        </Text>
      )}
      {bffAndGdsVersion && <Text as="span">{t`BFF & GDS version ${bffAndGdsVersion}`}</Text>}
    </Text>
  );
}

export default Version;
