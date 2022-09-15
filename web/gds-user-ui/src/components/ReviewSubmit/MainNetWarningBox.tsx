import { Trans } from '@lingui/react';
import WarningBox from 'components/WarningBox';
import { Text } from '@chakra-ui/react';

function MainNetWarningBox() {
  return (
    <WarningBox>
      <Text>
        <Trans id="If you would like to register for MainNet, please provide a MainNet Endpoint and Common Name.">
          If you would like to register for MainNet, please provide a MainNet Endpoint and Common
          Name.
        </Trans>
      </Text>
      <Text>
        <Trans id="Please note that TestNet and MainNet are separate networks that require different X.509 Identity Certificates.">
          Please note that TestNet and MainNet are separate networks that require different X.509
          Identity Certificates.
        </Trans>
      </Text>
    </WarningBox>
  );
}

export default MainNetWarningBox;
