import { FaChevronDown } from 'react-icons/fa';

const ArrowIcon = ({ open }: { open: boolean }) => (
  <FaChevronDown
    style={{
      transform: open ? 'rotate(180deg)' : undefined,
      transition: '200ms'
    }}
  />
);

export default ArrowIcon;
