import { Box, Heading, HStack, VStack } from '@chakra-ui/react';
import Card, { CardBody } from 'components/Card';
import TestNetCertificateProgressBar from 'components/testnetProgress/TestNetCertificateProgressBar.component';
import FormButton from 'components/ui/FormButton';

type CertificateLayoutProps = {
  children: React.ReactNode;
};

const CertificateLayout: React.FC<CertificateLayoutProps> = ({ children }) => {
  return (
    <>
      <Heading size="lg" mb="24px">
        Certificate Registration
      </Heading>
      <VStack spacing={3}>
        <Card maxW="100%">
          <CardBody>
            This multi-section form is an important step in the registration and certificate
            issuance process. The information you provide will be used to verify the legal entity
            that you represent and, where appropriate, will be available to verified TRISA members
            to facilitate compliance decisions. To assist in completing the registration form, the
            form is broken into multiple sections. No information is sent until you complete Section
            6 - Review & Submit.
          </CardBody>
        </Card>
        <Box width={'100%'}>
          <TestNetCertificateProgressBar />
        </Box>
        <Box pt="27px" w="100%">
          {children}
        </Box>
        <HStack width="100%" justifyContent="space-between" pt={5}>
          <FormButton borderRadius={5} w="100%" maxW="13rem">
            Previous
          </FormButton>
          <FormButton
            borderRadius={5}
            w="100%"
            maxW="13rem"
            backgroundColor="#FF7A59"
            _hover={{ backgroundColor: '#f07253' }}>
            Save & Next
          </FormButton>
          <FormButton borderRadius={5} w="100%" maxW="13rem">
            Save & Continue Later
          </FormButton>
        </HStack>
      </VStack>
    </>
  );
};

export default CertificateLayout;
