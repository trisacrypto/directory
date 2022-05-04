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

import { GoogleIcon } from 'components/Icon';

import { colors } from 'utils/theme';
import { useForm } from 'react-hook-form';

interface LoginProps {
  handleSignWithSocial: (event: React.FormEvent, type: string) => void;
  handleSignWithEmail: (data: any) => void;
  isLoading?: boolean;
  isError?: any;
}
interface IFormInputs {
  username: string;
  password: string;
}
const Login: React.FC<LoginProps> = (props) => {
  const {
    register,
    handleSubmit,
    formState: { errors }
  } = useForm<IFormInputs>();
  return (
    <Flex
      minWidth={'100vw'}
      align={'center'}
      justify={'center'}
      fontFamily={colors.font}
      fontSize={'xl'}
      marginTop={'10vh'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={8} mx={'auto'} maxW={'lg'} py={12} px={6} width={'100%'}>
        <Stack align={'left'}>
          <Heading fontSize={'xl'}>Log into your TRISA account.</Heading>
        </Stack>
        <Stack align={'center'} justify={'center'} fontFamily={colors.font}>
          <Button
            bg={'gray.100'}
            w="100%"
            height={'64px'}
            onClick={(event) => props.handleSignWithSocial(event, 'google')}
            color={'gray.600'}
            _hover={{
              background: useColorModeValue('gray.200', 'black'),
              color: useColorModeValue('gray.600', 'white')
            }}
            _focus={{
              borderColor: 'transparent'
            }}>
            <GoogleIcon h={24} />
            <Text as={'span'} ml={3}>
              Continue with google
            </Text>
          </Button>
          <Text py={3}>or</Text>
        </Stack>

        <Box
          rounded={'lg'}
          bg={useColorModeValue('white', 'transparent')}
          position={'relative'}
          bottom={5}>
          <form onSubmit={handleSubmit(props.handleSignWithEmail)}>
            <Stack spacing={4}>
              <FormControl id="email">
                <Input
                  type="email"
                  {...register('username')}
                  height={'64px'}
                  placeholder="Email Address"
                />
              </FormControl>
              <FormControl id="password">
                <Input
                  type="password"
                  {...register('password')}
                  height={'64px'}
                  placeholder="Password"
                />
              </FormControl>
              <Stack spacing={8} direction={['column', 'row']} py="10">
                <Button
                  bg={colors.system.blue}
                  color={'white'}
                  height={'57px'}
                  w={['full', '50%']}
                  type="submit"
                  _hover={{
                    bg: '#10aaed'
                  }}
                  _focus={{
                    borderColor: 'transparent'
                  }}>
                  Log In
                </Button>

                <Text lineHeight="57px">
                  {' '}
                  <Link href="/auth/forget"> Forgot password? </Link>
                </Text>
              </Stack>
            </Stack>
          </form>
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
};
export default Login;
