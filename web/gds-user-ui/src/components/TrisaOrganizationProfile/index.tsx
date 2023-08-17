import OrganizationalDetail from 'components/OrganizationProfile/OrganizationalDetail';
import TrisaImplementation from 'components/OrganizationProfile/TrisaImplementation';

import { handleError } from 'utils/utils';
import { useFetchCertificateStep } from 'hooks/useFetchCertificateStep';
import { StepEnum } from 'types/enums';

interface TrisaOrganizationProfileProps {
  type: 'mainnet' | 'testnet';
  status: string;
}
function TrisaOrganizationProfile({ type, status }: TrisaOrganizationProfileProps) {
  const { certificateStep, isFetchingCertificateStep, error } = useFetchCertificateStep({
    key: StepEnum.ALL
  });

  if (isFetchingCertificateStep) {
    return <>loading...</>;
  }

  if (error) {
    handleError(error);
    return null;
  }

  return (
    <div>
      <OrganizationalDetail data={certificateStep?.form} network={type} status={status} />
      <TrisaImplementation
        data={{
          mainnet: certificateStep?.form?.mainnet || {},
          testnet: certificateStep?.form?.testnet || {}
        }}
      />
    </div>
  );
}

export default TrisaOrganizationProfile;
