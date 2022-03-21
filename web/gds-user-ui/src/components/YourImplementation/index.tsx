import Card from "components/Card";

const YourImplementation: React.FC = () => {
  return (
    <Card borderRadius={10} maxW="100%">
      <Card.CardHeader size="sm">Your TRISA Implementation</Card.CardHeader>
      <Card.CardBody>
        Since TRISA is an open source, peer-to-peer Travel Rule solution, VASPs
        can set up and maintain their own TRISA server to exhange encrypted
        Travel Rule compliance data. At the same time, TRISA is designed to be
        interoperable. There are several Travel Rule solutions providers
        available on the market. If you are a customer, work with your Travel
        Ruie provider to integrate TRISA into your Travel Rule compliance
        workflow.
      </Card.CardBody>
    </Card>
  );
};

export default YourImplementation;
