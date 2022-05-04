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
interface AlertMessageProps {
  message?: string;
  status?: any;
  title?: string;
  hasBackBtn?: boolean;
}
const AlertMessage: React.FC<AlertMessageProps> = ({ status, title, message, hasBackBtn }) => {
  return (
    <Flex
      align={'center'}
      justify={'center'}
      fontFamily={colors.font}
      fontSize={'xl'}
      width={'100%'}
      mt={'10%'}
      bg={useColorModeValue('white', 'gray.800')}>
      <Stack width={'50%'}>
        <Alert
          status={status}
          variant="subtle"
          flexDirection="column"
          alignItems="center"
          justifyContent="center"
          textAlign="center">
          <AlertIcon boxSize="40px" mr={0} />
          <AlertTitle fontSize="lg" my={2}>
            {title ? title : 'Unable to Verify Contact'}
          </AlertTitle>
          <AlertDescription maxWidth="sm" py={4}>
            {message}
          </AlertDescription>
          {hasBackBtn && (
            <Box py={4}>
              <Button as={'a'} href={'/'}>
                Return to Directory
              </Button>
            </Box>
          )}
        </Alert>
      </Stack>
    </Flex>
  );
};
export default AlertMessage;
