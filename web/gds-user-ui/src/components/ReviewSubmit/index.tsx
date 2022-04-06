import { Button, Heading, VStack, Stack, Text } from '@chakra-ui/react';
import FormLayout from 'layouts/FormLayout';
interface ReviewSubmitProps {
  onSubmitHandler: (e: React.FormEvent, network: string) => void;
  isTestNetSent?: boolean;
  isMainNetSent?: boolean;
}
const ReviewSubmit: React.FC<ReviewSubmitProps> = ({
  onSubmitHandler,
  isTestNetSent,
  isMainNetSent
}) => {
  return (
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
      </Stack>
      <Stack spacing={10}>
        {isTestNetSent && (
          <FormLayout>
            <Text>
              Your{' '}
              <Text as="span" fontWeight="bold">
                TestNet
              </Text>{' '}
              registration form has been successfully submitted. You will receive a confirmation
              email from admin@trisa.io. In the email, you will receive instructions on next
            </Text>
            steps.
          </FormLayout>
        )}
        {isMainNetSent && (
          <FormLayout>
            <Text>
              Your{' '}
              <Text as="span" fontWeight="bold">
                MainNet
              </Text>{' '}
              registration form has been successfully submitted. You will receive a confirmation
              email from admin@trisa.io. In the email, you will receive instructions on next
            </Text>
            steps.
          </FormLayout>
        )}
      </Stack>
    </VStack>
  );
};

export default ReviewSubmit;
