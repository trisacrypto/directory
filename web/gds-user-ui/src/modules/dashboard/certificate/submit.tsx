import { SimpleDashboardLayout } from 'layouts';
import { Box, Heading, HStack, VStack } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import useCertificateStepper from 'hooks/useCertificateStepper';
import { RootStateOrAny, useSelector } from 'react-redux';
import ReviewSubmit from 'components/ReviewSubmit';
import CertificateStepLabel from 'components/testnetProgress/CertificateStepLabel';
const Certificate: React.FC = () => {
  const { nextStep, previousStep } = useCertificateStepper();
  const currentStep: number = useSelector((state: RootStateOrAny) => state.stepper.currentStep);
  const lastStep: number = useSelector((state: RootStateOrAny) => state.stepper.lastStep);

  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    await null;
    console.log('handleSubmitRegister', network);
  };
  return (
    // <DashboardLayout>
    // </DashboardLayout>
    <SimpleDashboardLayout>
      <CertificateStepLabel />
      <ReviewSubmit onSubmitHandler={handleSubmitRegister} />
    </SimpleDashboardLayout>
  );
};

export default Certificate;
