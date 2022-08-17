import { Link, ListItem, UnorderedList } from '@chakra-ui/react';
import { Trans } from '@lingui/react';
import Card from 'components/Card/Card';

const OpenSourceResources: React.FC = () => {
  return (
    <Card maxW="100%">
      <Card.CardHeader>
        <Trans id="Open Source Resources">Open Source Resources</Trans>
      </Card.CardHeader>
      <Card.CardBody>
        <UnorderedList>
          <ListItem color="#1F4CED">
            <Link href="https://github.com/trisacrypto/trisa" isExternal>
              <Trans id="TRISA Githubs repo">TRISA Github&apos;s repo</Trans>
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link href="https://trisa.dev/" isExternal>
              <Trans id="Documentation">Documentation</Trans>
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link
              href="https://github.com/trisacrypto/trisa/commit/436fd73fc48973ce09ccbae4260df6213d0c2894"
              isExternal>
              <Trans id="Reference implementation">Reference implementation</Trans>
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link isExternal>
              <Trans id="Meet Alice VASP, Bob VASP and Evil VASP">
                Meet Alice VASP, Bob VASP and &quot;Evil&quot; VASP
              </Trans>
            </Link>
          </ListItem>
        </UnorderedList>
      </Card.CardBody>
    </Card>
  );
};

export default OpenSourceResources;
