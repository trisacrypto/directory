import { Box, Heading, VStack } from '@chakra-ui/react';
import BasicDetails from 'components/BasicDetail';
import Card, { CardBody } from 'components/Card';
import TestNetCertificateProgressBar from 'components/testnetProgress/TestNetCertificateProgressBar.component';
import DashboardLayout from 'layouts/DashboardLayout';
import CertificateLayout from './CertificateLayout';

const Certificate: React.FC = () => {
  return (
    <DashboardLayout>
      <CertificateLayout>
        <BasicDetails />
      </CertificateLayout>
    </DashboardLayout>
  );
};

export default Certificate;
