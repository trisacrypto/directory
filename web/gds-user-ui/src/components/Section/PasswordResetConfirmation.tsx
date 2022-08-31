import { Flex, Stack, Heading, Text, useColorModeValue } from '@chakra-ui/react';
import { Trans } from '@lingui/react';

export default function PasswordResetConfirmation(props: any) {
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
            <Trans id="Thank you. We have sent instructions to reset your password to">
              Thank you. We have sent instructions to reset your password to
            </Trans>
            <Text as={'span'} fontWeight={'bold'}>
              {props.email}
            </Text>
            .{' '}
            <Trans id="The link to reset your password expires in 24 hours.">
              The link to reset your password expires in 24 hours.
            </Trans>
          </Heading>
        </Stack>
      </Stack>
    </Flex>
  );
}
