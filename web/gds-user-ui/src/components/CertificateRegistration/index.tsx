import { Stack, Heading } from '@chakra-ui/react';
import Card from 'components/ui/Card';
import { Trans } from '@lingui/react';
import CertificateRegistrationTable from './CertificateRegistrationTable';

const CertificateRegistration = ({ data }: any) => (
  <Card>
    <Stack p={4} mb={5}>
      <Heading fontSize="20px" fontWeight="bold" pb=".5rem">
        <Trans id="Certificate Registration Process">Certificate Registration Process</Trans>
      </Heading>
    </Stack>
    <CertificateRegistrationTable />
  </Card>
);

export default CertificateRegistration;
