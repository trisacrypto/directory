import { Text } from "@chakra-ui/react";
import { Trans } from "@lingui/macro";
import Card from "components/ui/Card";
import React from 'react';

const DirectoryNotification = () => {
    return (
        <Card maxW="100%" marginBottom={6}>
        <Card.Body>
          <Text>
            <Trans>
              The TRISA Member Directory is for informational purposes and is meant to foster collaboration between members. It is available to verified TRISA contacts only. If you cannot view the list, please complete the registration process. You can view the member list after verification is complete.
            </Trans>
          </Text>
        </Card.Body>
      </Card>
    );
};

export default React.memo(DirectoryNotification);
