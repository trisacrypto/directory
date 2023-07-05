import { act, render, screen } from 'utils/test-utils';
import MemberTableRows, { MemberTableRowsProps } from '../MemberTableRows';

import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { mainnetMembersMockValue } from '../../__mocks__';
import { Table, Tbody } from '@chakra-ui/react';

function renderComponent(rowProps?: MemberTableRowsProps) {
  const defaultProps: MemberTableRowsProps = {
    rows: mainnetMembersMockValue.vasps,
    hasError: false,
    isLoading: false
  };
  const props = rowProps || defaultProps;
  return render(
    <Table>
      <Tbody>
        <MemberTableRows {...props} />
      </Tbody>
    </Table>
  );
}

describe('MemberTable', () => {
  beforeAll(() => {
    act(() => {
      dynamicActivate('en');
    });
  });

  it('should render', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render the mainnet data in the table rows', () => {
    renderComponent();
    const memberName1 = screen.getByText('FireCoin Exchange');
    const joined1 = screen.getByText('Apr 20, 2022');
    const lastUpdated1 = screen.getByText('Nov 19, 2022');
    const memberName2 = screen.getByText('Test Machine');
    const joined2 = screen.getByText('Mar 22, 2023');
    const lastUpdated2 = screen.getByText('May 10, 2023');
    const memberName3 = screen.getByText('New Coin Exchange');
    const joined3 = screen.getByText('Jan 20, 2023');
    const lastUpdated3 = screen.getByText('Jun 18, 2023');
    expect(memberName1).toBeInTheDocument();
    expect(joined1).toBeInTheDocument();
    expect(lastUpdated1).toBeInTheDocument();
    expect(memberName2).toBeInTheDocument();
    expect(joined2).toBeInTheDocument();
    expect(lastUpdated2).toBeInTheDocument();
    expect(memberName3).toBeInTheDocument();
    expect(joined3).toBeInTheDocument();
    expect(lastUpdated3).toBeInTheDocument();
  });

  it('should render the unverified member error', () => {
    const propsMock = {
      rows: [],
      hasError: true,
      isLoading: false
    };
    renderComponent(propsMock);
    const unverifiedMember = screen.getByTestId('unverified-member');
    expect(unverifiedMember).toBeInTheDocument();
  });
});
