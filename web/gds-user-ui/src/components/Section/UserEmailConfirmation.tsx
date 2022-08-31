import {
  Flex,
  Stack,
  Button,
  useColorModeValue,
} from '@chakra-ui/react';
import { t } from '@lingui/macro';
import { Trans } from '@lingui/react';

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
          <AlertMessage message={props.message} status="success" title={t`Contact Verified`} />
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
            <Trans id="Log In">Log In</Trans>
          </Button>
        </Stack>
      </Stack>
    </Flex>
  );
};
export default AuthEmailConfirmation;
