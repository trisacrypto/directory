import { useState, useEffect } from 'react';
import { Button, Heading, VStack, Stack, Text, useDisclosure } from '@chakra-ui/react';
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

  useEffect(() => {
    if (isSent) {
      onOpen();
    }
  }, [isSent]);
  return (
    <>
      <VStack align="start" mt="2rem">
        <Heading size="md">Registration Submission</Heading>
        <FormLayout>
          <Text>
            You must submit your registration for TestNet and MainNet separately.{' '}
            <Text as="span" fontWeight="bold">
              Note:
            </Text>{' '}
            You will receive two separate emails with confirmation links for each registration. You
            must click on each confirmation link to complete the registration process.
            <Text as="span" fontWeight="bold">
              Failure to click either confirmation will result in an incomplete registration.
            </Text>
          </Text>
        </FormLayout>
        <Stack direction={['column', 'row']} justifyContent="space-around" py={14} width="100%">
          <Button
            bgColor="#ff7a59f0"
            color="#fff"
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
          <Button
            bgColor="#23a7e0e8"
            color="#fff"
            size="lg"
            py="2.5rem"
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
        </Stack>
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
    </>
  );
};

export default ReviewSubmit;
