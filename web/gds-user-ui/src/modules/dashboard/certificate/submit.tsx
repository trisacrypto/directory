import { SimpleDashboardLayout } from 'layouts';
import CertificateStepLabel from 'components/TestnetProgress/CertificateStepLabel';
import { lazy, Suspense } from 'react';
import Loader from 'components/Loader';

const ReviewSubmit = lazy(() => import('components/ReviewSubmit'));

const Certificate: React.FC = () => {
  const handleSubmitRegister = async (event: React.FormEvent, network: string) => {
    event.preventDefault();
    await null;
    console.log('handleSubmitRegister', network);
  };
  return (
    <SimpleDashboardLayout>
      <Suspense fallback={<Loader />}>
        <CertificateStepLabel />
        <ReviewSubmit onSubmitHandler={handleSubmitRegister} />
      </Suspense>
    </SimpleDashboardLayout>
  );
};

export default Certificate;
