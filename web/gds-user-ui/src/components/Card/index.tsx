import { Box, BoxProps, Heading, HeadingProps, Stack } from "@chakra-ui/react";

const CardHeader: React.FC<HeadingProps> = ({ children, ...props }) => {
  return (
    <Heading size="sm" {...props}>
      {children}
    </Heading>
  );
};

export const CardBody: React.FC<BoxProps> = (props) => {
  return <Box {...props} />;
};

interface CardProps {
  CardHeader: React.FC<HeadingProps>;
  CardBody: React.FC<BoxProps>;
}

const Card: React.FC<BoxProps> & CardProps = (props) => {
  return (
    <Stack
      border="1px solid #C4C4C4"
      padding="20px"
      w="100%"
      maxW="300px"
      spacing="17px"
      borderRadius={10}
      {...props}
    />
  );
};

Card.CardHeader = CardHeader;
Card.CardBody = CardBody;

export default Card;
