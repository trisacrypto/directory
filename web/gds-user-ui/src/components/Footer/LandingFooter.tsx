/* eslint-disable no-console */
import React, { useEffect, useState } from 'react';
import { Flex, Text, Link, useColorModeValue } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import useAxios from 'hooks/useAxios';
import {
  getAppVersionNumber,
  getBffAndGdsVersion,
  getAppGitVersion,
  isProdEnv
} from 'application/config';
import { Trans } from '@lingui/react';

const Footer = (): React.ReactElement => {
  const [appVersion, setAppVersion] = useState<any>();
  const [gitRevision, setGitRevision] = useState<any>();
  const [bffAndGdsVersion, setBffAndGdsVersion] = useState<any>();
  const fetchAsyncBffAndGdsVersion = async () => {
    const request = await getBffAndGdsVersion();
    if (request) {
      setBffAndGdsVersion(request.version);
    }
  };

  useEffect(() => {
    // console.log(data);
    const getAppVersion = getAppVersionNumber();
    const getGitRevision = getAppGitVersion();
    setAppVersion(getAppVersion);
    setGitRevision(getGitRevision);
    fetchAsyncBffAndGdsVersion();
  }, []);
  //  log this out in the console
  if (isProdEnv) {
    console.log('appVersion', appVersion);
    console.log('gitRevision', gitRevision);
    console.log('bffAndGdsVersion', bffAndGdsVersion);
  }
  return (
    <Flex
      bg={useColorModeValue(colors.system.gray, 'transparent')}
      color="white"
      width="100%"
      justifyContent="center"
      alignItems="center"
      direction="column"
      padding={4}
      // position={'absolute'}
      bottom={0}>
      <Flex width="100%" wrap="wrap">
        <Text width="100%" textAlign="center" color="white" fontSize="sm">
          <Trans id="A component of">A component of</Trans>{' '}
          <Link href="https://trisa.io" color={colors.system.cyan}>
            <Trans id="the TRISA architecture">the TRISA architecture</Trans>
          </Link>{' '}
          <Trans id="for Cryptocurrency Travel Rule compliance.">
            for Cryptocurrency Travel Rule compliance.
          </Trans>
        </Text>
        <Text width="100%" textAlign="center" color="white" fontSize="sm">
          <Trans id="Created and maintained by">Created and maintained by</Trans>{' '}
          <Link href="https://rotational.io" color={colors.system.cyan}>
            {' '}
            Rotational Labs
          </Link>{' '}
          <Trans id="in partnership with">in partnership with</Trans>{' '}
          <Link href="https://cyphertrace.com" color={colors.system.cyan}>
            {' '}
            CipherTrace
          </Link>{' '}
          on behalf of{' '}
          <Link href="https://trisa.io" color={colors.system.cyan}>
            TRISA
          </Link>{' '}
        </Text>

        {isProdEnv && (
          <Text width="100%" textAlign="center" color="white" fontSize="12" pt={1}>
            {appVersion && <Text as="span">App version {appVersion} - </Text>}
            {gitRevision && <Text as="span">Git Revision {gitRevision} - </Text>}
            {bffAndGdsVersion && <Text as="span">BFF & GDS version {bffAndGdsVersion} </Text>}
          </Text>
        )}
      </Flex>
    </Flex>
  );
};

export default Footer;
