import { Link, ListItem, UnorderedList } from '@chakra-ui/react';
import Card from 'components/Card/Card';

const OpenSourceResources: React.FC = () => {
  return (
    <Card maxW="100%">
      <Card.CardHeader>Open Source Resources</Card.CardHeader>
      <Card.CardBody>
        <UnorderedList>
          <ListItem color="#1F4CED">
            <Link href="https://github.com/trisacrypto/trisa" isExternal>
              TRISA Github&apos;s repo
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link href="https://trisa.dev/" isExternal>
              Documentation
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link
              href="https://github.com/trisacrypto/trisa/commit/436fd73fc48973ce09ccbae4260df6213d0c2894"
              isExternal>
              Reference implementation
            </Link>
          </ListItem>
          <ListItem color="#1F4CED">
            <Link isExternal>Meet Alice VASP, Bob VASP and &quot;Evil&quot; VASP</Link>
          </ListItem>
        </UnorderedList>
      </Card.CardBody>
    </Card>
  );
};

export default OpenSourceResources;
