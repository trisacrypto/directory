import { useState, useEffect } from 'react';
import {
  Button,
  Heading,
  VStack,
  Stack,
  Text,
  useDisclosure,
  Box,
  Flex,
  Link
} from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
import ConfirmationModal from 'components/ReviewSubmit/ConfirmationModal';
interface ReviewSubmitProps {
  onSubmitHandler: (e: React.FormEvent, network: string) => void;
  isTestNetSent?: boolean;
  isMainNetSent?: boolean;
  result?: any;
}
const ReviewSubmit: React.FC<ReviewSubmitProps> = ({
  onSubmitHandler,
  isTestNetSent,
  isMainNetSent,
  result
}) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  const isSent = isTestNetSent || isMainNetSent;
  const [testnet, setTestnet] = useState(false);
  const [mainnet, setMainnet] = useState(false);
  const getTestnetFromLocalStorage = localStorage.getItem('isTestNetSent');
  const getMainnetFromLocalStorage = localStorage.getItem('isMainNetSent');
  useEffect(() => {
    if (getTestnetFromLocalStorage === 'true') {
      setTestnet(true);
    }
    if (getMainnetFromLocalStorage === 'true') {
      setMainnet(true);
    }
  }, [getTestnetFromLocalStorage, getMainnetFromLocalStorage]);
  useEffect(() => {
    if (isSent) {
      onOpen();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isTestNetSent, isMainNetSent]);
  return (
    <>
      <Flex>
        <VStack mt="2rem">
          <Stack align="start" w="full">
            <Heading size="md" pt={2}>
              Registration Submission
            </Heading>
          </Stack>

          <FormLayout>
            <Text>
              You must submit your registration for TestNet and MainNet separately.{' '}
              <Text as="span" fontWeight="bold">
                Note:
              </Text>{' '}
              You will receive two separate emails with confirmation links for each registration.
              You must click on each confirmation link to complete the registration process.
              <Text as="span" fontWeight="bold">
                Failure to click either confirmation will result in an incomplete registration.
              </Text>
            </Text>
          </FormLayout>
          <Stack
            direction={['column', 'row']}
            justifyContent="space-around"
            py={14}
            width="100%"
            spacing={10}>
            <Stack bg={'white'}>
              <Stack px={6} mb={5} pt={4}>
                <Heading size="sm" mt={2}>
                  TESTNET SUBMISSION
                </Heading>
                <Text>
                  Click below to submit your{' '}
                  <Text as="span" fontWeight={'bold'}>
                    TestNet
                  </Text>{' '}
                  registration. Upon submission, you will receive an email with a confirmation link.
                  You must click the confirmation link to complete the registration process. Failure
                  to click the confirmation link will result in an incomplete registration.
                </Text>
                <Text>
                  A physical verification check in the form of a phone call{' '}
                  <Text as="span" fontWeight={'bold'}>
                    is not required
                  </Text>{' '}
                  for TestNet registration so your TestNet certificate will be issued upon review by
                  the validation team.
                </Text>
                <Text>
                  If you would like to edit your registration form before submitting, please return
                  to the{' '}
                  <Link color={'blue'} href="/certificate/registration" fontWeight={'bold'}>
                    Review page
                  </Link>
                  .
                </Text>
              </Stack>
              <Stack
                alignContent={'center'}
                justifyContent={'center'}
                mx="auto"
                pb={4}
                alignItems={'center'}>
                <Button
                  bgColor="#ff7a59f0"
                  color="#fff"
                  isDisabled={testnet}
                  size="lg"
                  py="2.5rem"
                  whiteSpace="normal"
                  maxW="200px"
                  width="100%"
                  boxShadow="lg"
                  onClick={(e) => {
                    onSubmitHandler(e, 'testnet');
                  }}
                  _hover={{
                    bgColor: '#f55c35'
                  }}>
                  Submit TestNet Registration
                </Button>
              </Stack>
            </Stack>

            <Stack bg={'white'}>
              <Stack px={6} mb={5} pt={4}>
                <Heading size="sm" mt={2}>
                  MAINNET SUBMISSION
                </Heading>
                <Text>
                  Click below to submit your{' '}
                  <Text as="span" fontWeight={'bold'}>
                    MainNet
                  </Text>{' '}
                  registration. Upon submission, you will receive an email with a confirmation link.
                  You must click the confirmation link to complete the registration process. Failure
                  to click the confirmation link will result in an incomplete registration.
                </Text>
                <Text>
                  A physical verification check in the form of a phone call{' '}
                  <Text as="span" fontWeight={'bold'}>
                    is required
                  </Text>{' '}
                  for MainNet registration so your TestNet certificate will be issued upon review by
                  the validation team.
                </Text>
                <Text>
                  If you would like to edit your registration form before submitting, please return
                  to the{' '}
                  <Link color={'blue'} href="/certificate/registration" fontWeight={'bold'}>
                    Review page
                  </Link>
                </Text>
              </Stack>
              <Stack
                alignContent={'center'}
                justifyContent={'center'}
                mx="auto"
                alignItems={'center'}
                pb={4}>
                <Button
                  bgColor="#23a7e0e8"
                  color="#fff"
                  size="lg"
                  py="2.5rem"
                  isDisabled={mainnet}
                  whiteSpace="normal"
                  boxShadow="lg"
                  maxW="200px"
                  onClick={(e) => {
                    onSubmitHandler(e, 'mainnet');
                  }}
                  width="100%"
                  _hover={{
                    bgColor: '#189fda'
                  }}>
                  Submit MainNet Registration
                </Button>
              </Stack>
            </Stack>
          </Stack>

          <Box alignItems={'center'} textAlign="center" mx={'auto'}>
            <Button
              bgColor="#555151"
              color="#fff"
              as="a"
              href="/certificate/registration"
              size="lg"
              py="2.5rem"
              whiteSpace="normal"
              boxShadow="lg"
              maxW="200px"
              width="100%"
              _hover={{
                bgColor: '#555151'
              }}>
              Back to Review Page
            </Button>
          </Box>
        </VStack>
        {isSent && (
          <ConfirmationModal
            isOpen={isOpen}
            onClose={onClose}
            id={result?.id}
            pkcs12password={result?.pkcs12password}
            message={result?.message}
            status={result?.status}
            size={'xl'}
          />
        )}
      </Flex>
    </>
  );
};

export default ReviewSubmit;
