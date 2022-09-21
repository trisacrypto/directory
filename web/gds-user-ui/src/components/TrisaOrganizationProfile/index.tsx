import OrganizationalDetail from 'components/OrganizationProfile/OrganizationalDetail';
import TrisaImplementation from 'components/OrganizationProfile/TrisaImplementation';
import { getRegistrationDefaultValue } from 'modules/dashboard/registration/utils';
import { useAsync } from 'react-use';
import { handleError } from 'utils/utils';

function TrisaOrganizationProfile() {
  const { value, error, loading } = useAsync(getRegistrationDefaultValue);

  if (loading) {
    return <>loading...</>;
  }

  if (error) {
    handleError(error);
    return null;
  }

  return (
    <div>
      <OrganizationalDetail data={value} />
      <TrisaImplementation data={{ mainnet: value.mainnet, testnet: value.testnet }} />
    </div>
  );
}

export default TrisaOrganizationProfile;
