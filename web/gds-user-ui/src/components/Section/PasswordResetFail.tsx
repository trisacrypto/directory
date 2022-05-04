import {
  Flex,
  Box,
  FormControl,
  FormLabel,
  Input,
  Checkbox,
  Stack,
  Link,
  Button,
  Heading,
  Text,
  useColorModeValue,
  Image
} from '@chakra-ui/react';

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
            Sorry. We could not find a user account with the email address
            <Text as={'span'}>{props.email}</Text>.
          </Heading>
        </Stack>

        <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          position={'relative'}
          bottom={5}>
          <Text textAlign="center">
            Not a TRISA Member?{' '}
            <Link href="/register" color={colors.system.cyan}>
              {' '}
              Join the TRISA network today.{' '}
            </Link>
          </Text>
        </Box>
      </Stack>
    </Flex>
  );
}
