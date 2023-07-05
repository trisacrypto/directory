import { Heading } from '@chakra-ui/react';
import MemberTable from './MemberTable';
import { Trans } from '@lingui/macro';
import Loader from 'components/Loader';
import { useFetchMembers } from './hook/useFetchMembers';
// import { mainnetMembersMockValue } from './__mocks__';

const MemberPage: React.FC = () => {
  const { members, isFetchingMembers, error } = useFetchMembers();

  // const vasps = mainnetMembersMockValue.vasps;
  // console.log('vasps', vasps);

  if (isFetchingMembers) return <Loader />;
  console.log('members', members);
  return (
    <>
      <Heading marginBottom="69px">
        <Trans>TRISA Member Directory</Trans>
      </Heading>

      {error && <p>error </p>}
      { members && <MemberTable data={members} />}
    </>
  );
};

export default MemberPage;
