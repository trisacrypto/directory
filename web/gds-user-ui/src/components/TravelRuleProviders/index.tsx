import { Link, ListItem, UnorderedList } from '@chakra-ui/react';
import Card from 'components/Card/Card';

const TravelRuleProviders: React.FC = () => {
  return (
    <Card maxW="100%">
      <Card.CardHeader>3rd Party Travel Rule Providers</Card.CardHeader>
      <Card.CardBody>
        <UnorderedList>
          <ListItem color="#1F4CED">
            <Link href="https://ciphertrace.com/" isExternal>
              Cyphertrace
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link href="https://www.sygna.io/" isExternal>
              Synga Bridge
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link href="https://notabene.id/" isExternal>
              NotaBene{' '}
            </Link>
            (not interoperable)
          </ListItem>
          <ListItem color="#1F4CED">
            <Link href="https://openvasp.org/" isExternal>
              OpenVASP{' '}
            </Link>
            (not interoperable)
          </ListItem>
        </UnorderedList>
      </Card.CardBody>
    </Card>
  );
};

export default TravelRuleProviders;
