import { Flex, Stack, useColorModeValue } from '@chakra-ui/react';
import { ReactNode } from 'react';

type AuthLayoutProps = {
  children: ReactNode;
};

function AuthLayout({ children, ...rest }: AuthLayoutProps) {
  return (
    <Flex fontSize={'xl'} bg={useColorModeValue('white', 'gray.800')} {...rest}>
      <Stack spacing={8} mx="auto" w="100%" maxW="lg" py={12}>
        {children}
      </Stack>
    </Flex>
  );
}

export default AuthLayout;
