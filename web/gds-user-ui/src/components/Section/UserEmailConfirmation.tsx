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
  Image,
  Alert,
  AlertIcon,
  AlertTitle,
  AlertDescription
} from '@chakra-ui/react';

import { colors } from '../../utils/theme';
interface AuthEmailConfirmationProps {
  message?: string;
}
import AlertMessage from '../ui/AlertMessage';
const AuthEmailConfirmation: React.FC<AuthEmailConfirmationProps> = (props) => {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={'open sans'}
      fontWeight={'bold'}
      fontSize={'xl'}
      mt={'10%'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack spacing={6} mx={'auto'} maxW={'xl'} py={12} px={6}>
        <Stack>
          <AlertMessage message={props.message} status="success" title={'Contact Verified'} />
        </Stack>
        <Stack spacing={8} direction={['column', 'row']} py="10">
          <Button
            bg={colors.system.blue}
            color={'white'}
            height={'57px'}
            w={['full', '50%']}
            _hover={{
              bg: '#10aaed'
            }}
            _focus={{
              borderColor: 'transparent'
            }}>
            Log In
          </Button>
        </Stack>
      </Stack>
    </Flex>
  );
};
export default AuthEmailConfirmation;
