import React from 'react';
import { Flex, Text, Link, useColorModeValue, Stack } from '@chakra-ui/react';
import { colors } from 'utils/theme';
import { Trans } from '@lingui/react';
import Version from './Version';

const Footer = (): React.ReactElement => {
  return (
    <footer style={{ width: '100%' }}>
      <Stack
        bg={useColorModeValue(colors.system.gray, 'transparent')}
        color="white"
        width="100%"
        justifyContent="center"
        alignItems="center"
        direction="column"
        padding={4}>
        <Flex width="100%" wrap="wrap">
          <Text width="100%" textAlign="center" color="white" fontSize="sm">
            <Trans id="A component of">A component of</Trans>{' '}
            <Link href="https://trisa.dev" color={colors.system.cyan}>
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
            <Trans id="on behalf of">on behalf of</Trans>{' '}
            <Link href="https://trisa.dev" color={colors.system.cyan}>
              TRISA
            </Link>{' '}
          </Text>
          <Version />
        </Flex>
      </Stack>
    </footer>
  );
};

export default React.memo(Footer);
