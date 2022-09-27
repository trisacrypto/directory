import { Box, BoxProps, HeadingProps } from '@chakra-ui/react';
import CardBody from './CardBody';
import CardHeader from './CardHeader';

interface CardProps {
  CardHeader: React.FC<HeadingProps>;
  CardBody: React.FC<any>;
}

const Card: React.FC<BoxProps> & CardProps = (props) => {
  const { children, ...rest } = props;
  return (
    <Box
      border="2px solid #E5EDF1"
      borderRadius="10px"
      padding={{ base: 3, md: 9 }}
      fontFamily="Open Sans"
      data-testid="card"
      {...rest}>
      {children}
    </Box>
  );
};

Card.CardHeader = CardHeader;
Card.CardBody = CardBody;

export default Card;
