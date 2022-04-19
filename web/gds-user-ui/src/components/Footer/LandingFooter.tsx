import React from 'react';
import { Flex, Text, Link, useColorModeValue } from '@chakra-ui/react';
import { colors } from 'utils/theme';
const Footer = (): React.ReactElement => {
  return (
    <Flex
      bg={useColorModeValue(colors.system.gray, 'transparent')}
      color="white"
      width="100%"
      justifyContent="center"
      alignItems="center"
      direction="column"
      padding={4}
      position={'absolute'}
      bottom={0}>
      <Flex width="100%" wrap="wrap">
        <Text width="100%" textAlign="center" color="white" fontSize="sm">
          A component of{' '}
          <Link href="https://trisa.io" color={colors.system.cyan}>
            the TRISA architecture
          </Link>{' '}
          for Cryptocurrency Travel Rule compliance.
        </Text>
        <Text width="100%" textAlign="center" color="white" fontSize="sm">
          Created and maintained by{' '}
          <Link href="https://rotational.io" color={colors.system.cyan}>
            {' '}
            Rotational Labs
          </Link>{' '}
          in partnership with{' '}
          <Link href="https://cyphertrace.com" color={colors.system.cyan}>
            {' '}
            CipherTrace
          </Link>{' '}
          on behalf of{' '}
          <Link href="https://trisa.io" color={colors.system.cyan}>
            TRISA
          </Link>{' '}
          .
        </Text>
      </Flex>
    </Flex>
  );
};

export default Footer;
