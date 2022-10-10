import { Box, BoxProps, Heading, HeadingProps } from '@chakra-ui/react';

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
  Header: React.FC<HeadingProps>;
  Body: React.FC<any>;
}

const Card: React.FC<BoxProps> & CardProps = (props) => {
  const { children, ...rest } = props;
  return (
    <Box
      border="2px solid #E5EDF1"
      borderRadius="10px"
      padding={{ base: 3, md: 9 }}
      fontFamily="Open Sans"
      bg="white"
      {...rest}>
      {children}
    </Box>
  );
};

Card.Header = CardHeader;
Card.Body = CardBody;

export default Card;
