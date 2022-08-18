import { Flex, Box, Stack, Heading, Text, useColorModeValue } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import { Link } from 'react-router-dom';

import { colors } from '../../utils/theme';

interface PasswordResetProps {
  email: string;
}
export default function PasswordReset(props: PasswordResetProps) {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontSize={'xl'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={12} mx={'auto'} maxW={'lg'} py={12} px={6}>
        <Stack align={'center'}>
          <Heading fontSize={'xl'}>
            <Trans id="Sorry. We could not find a user account with the email address">
              Sorry. We could not find a user account with the email address
            </Trans>
            <Text as={'span'}>{props.email}</Text>.
          </Heading>
        </Stack>

        <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          position={'relative'}
          bottom={5}>
          <Text textAlign="center">
            <Trans id="Not a TRISA Member?">Not a TRISA Member?</Trans>{' '}
            <Link to="/register" color={colors.system.cyan}>
              <Trans id="Join the TRISA network today.">Join the TRISA network today.</Trans>
            </Link>
          </Text>
        </Box>
      </Stack>
    </Flex>
  );
}
