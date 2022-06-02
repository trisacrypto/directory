import { ListItem, UnorderedList } from '@chakra-ui/react';
import Card from 'components/Card/Card';

type MissingRequireProps = {
  missingFields: Record<string, string | number | null>;
};

const MissingRequire: React.FC<MissingRequireProps> = ({ missingFields }) => {
  return (
    <>
      <Card borderWidth="2px" borderStyle="solid" borderColor="red.500" color="red.500">
        <Card.CardHeader>Please complete all details</Card.CardHeader>
        <Card.CardBody>
          <UnorderedList mt={2}>
            {Object.entries(missingFields).map(([k, v], idx) => (
              <ListItem key={idx}>{v}</ListItem>
            ))}
          </UnorderedList>
        </Card.CardBody>
      </Card>
    </>
  );
};

export default MissingRequire;
