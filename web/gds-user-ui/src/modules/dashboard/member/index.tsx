import { Heading } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
import { useFetchMembers } from './hook/useFetchMembers';
import DirectoryNotification from './Components/DirectoryNotification';

const MemberPage: React.FC = () => {
  const { members, isFetchingMembers, error } = useFetchMembers();

  if (isFetchingMembers) return <Loader />;
  console.log('members', members);
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>
      <DirectoryNotification />

      {error && <p>error </p>}
      {members && <MemberTable data={members} />}
    </>
  );
};

export default MemberPage;
