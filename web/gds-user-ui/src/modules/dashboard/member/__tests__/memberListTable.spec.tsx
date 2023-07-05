import { act, fireEvent, render, screen } from 'utils/test-utils';
import MemberTable from "../MemberTable";
import { dynamicActivate } from 'utils/i18nLoaderHelper';
import { mainnetMembersMockValue, testnetMembersMockValue } from '../__mocks__';

const mainnetData = mainnetMembersMockValue.vasps;
const testnetData = testnetMembersMockValue.vasps;

function renderComponent() {
  return render(<MemberTable data={mainnetData} />);
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
      })

        it('should render the correct table headers', () => {
            renderComponent();
            const name = screen.getByTestId('name-header');
            const joined = screen.getByTestId('joined-header');
            const lastUpdated = screen.getByTestId('last-updated-header');
            const network = screen.getByTestId('network-header');
            const status = screen.getByTestId('status-header');
            const actions = screen.getByTestId('actions-header');
            expect(name).toBeInTheDocument();
            expect(joined).toBeInTheDocument();
            expect(lastUpdated).toBeInTheDocument();
            expect(network).toBeInTheDocument();
            expect(status).toBeInTheDocument();
            expect(actions).toBeInTheDocument();
        });
        
        it('should render the mainnet data in the table rows by default', () => {
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

        it('should render testnet data if testnet is selected in the select network component', () => {
          render(<MemberTable data={testnetData} />)
          const selectNetwork = screen.getByTestId('select-network');
          fireEvent.change(selectNetwork, { target: { value: 'testnet' } });
          expect(selectNetwork).toBeInTheDocument();
          expect(selectNetwork).toHaveValue('testnet');
          const memberName1 = screen.getByText('SendCoin VASP');
          const joined1 = screen.getByText('Feb 10, 2022');
          const lastUpdated1 = screen.getByText('Apr 9, 2023');
          const memberName2 = screen.getByText('Example Crypto');
          const joined2 = screen.getByText('Dec 1, 2021');
          const lastUpdated2 = screen.getByText('Feb 23, 2023');
          const memberName3 = screen.getByText('SpudCoin');
          const joined3 = screen.getByText('Jul 23, 2021');
          const lastUpdated3 = screen.getByText('Nov 27, 2022');
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
});
