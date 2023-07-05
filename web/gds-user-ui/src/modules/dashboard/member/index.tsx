import { Heading } from '@chakra-ui/react';

import MemberTable from './components/MemberTable';
import { Trans } from '@lingui/macro';
import DirectoryNotification from './components/DirectoryNotification';
import FormLayout from 'layouts/FormLayout';
import MemberSelectNetwork from './components/memberNetworkSelect';
import MemberTableHeader from './components/MemberTableHeader';
const MemberPage: React.FC = () => {
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>
      <DirectoryNotification />
      <FormLayout overflowX={'scroll'}>
        <MemberTableHeader />
        <MemberSelectNetwork />
        <MemberTable />
      </FormLayout>
    </>
  );
};

export default MemberPage;
