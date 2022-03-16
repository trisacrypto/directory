import { Flex, Stack, Text, useColorModeValue } from '@chakra-ui/react';

export default function AuthEmailConfirmation() {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontSize={'xl'}
      mt={'10%'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={6} mx={'auto'} maxW={'xl'} py={12} px={6}>
        <Stack align={'center'}>
          <Text color={useColorModeValue('gray.600', 'white')}>
            <Text as={'span'} fontWeight={'bold'}>
              {' '}
              Thank you for creating your TRISA account.{' '}
            </Text>{' '}
            <br />
            You must verify your email address. An email with verification instructions has been
            sent to your email address. Please complete the email verification process in 24 hours.
            The email verification link will expire in 24 hours.
          </Text>
        </Stack>
      </Stack>
    </Flex>
  );
}
