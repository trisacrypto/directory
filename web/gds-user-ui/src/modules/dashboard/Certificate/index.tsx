import { SimpleDashboardLayout } from 'layouts';
import { Box, Heading, HStack, VStack } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import TestNetCertificateProgressBar from 'components/TestnetProgress/TestNetCertificateProgressBar.component';
import FormButton from 'components/ui/FormButton';
import CertificateLayout from 'layouts/CertificateLayout';
import useCertificateStepper from 'hooks/useCertificateStepper';

const Certificate: React.FC = () => {
  const { nextStep, previousStep } = useCertificateStepper();

  return (
    // <DashboardLayout>
    //   <CertificateLayout>
    //     <BasicDetails />
    //   </CertificateLayout>
    // </DashboardLayout>
    <SimpleDashboardLayout>
      <>
        <Heading size="lg" mb="24px">
          Certificate Registration
        </Heading>
        <VStack spacing={3}>
          <Card maxW="100%">
            <Card.Body>
              This multi-section form is an important step in the registration and certificate
              issuance process. The information you provide will be used to verify the legal entity
              that you represent and, where appropriate, will be available to verified TRISA members
              to facilitate compliance decisions. To assist in completing the registration form, the
              form is broken into multiple sections. No information is sent until you complete
              Section 6 - Review & Submit.
            </Card.Body>
          </Card>
          <Box width={'100%'}>
            <TestNetCertificateProgressBar />
          </Box>

          <HStack width="100%" justifyContent="space-between" pt={5}>
            <FormButton onClick={() => previousStep()} borderRadius={5} w="100%" maxW="13rem">
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
            <FormButton onClick={() => nextStep()} borderRadius={5} w="100%" maxW="13rem">
              Save & Continue Later
            </FormButton>
          </HStack>
        </VStack>
      </>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
