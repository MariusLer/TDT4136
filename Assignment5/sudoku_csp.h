#include <vector>
#include <string>
#include <iostream>
#include <map>
#include <utility>

using namespace std;

/*
typedef vector<string> varriable_set
typedef map<string, vector<string>> legal_values_map
typedef map<string, map<string, vector<pair<string, string>>>> constraint_map
*/

class CSP{
private:
  vector<string> variables;
  map<string, vector<string>> domain;
  map<string, map<string, vector<pair<string, string>>>> constraint;

public:
  void add_variable(string name, vector<string> domain);
  vector<pair<string, string>> get_all_possible_pairs(vector<string> a, vector<string> b);
  vector<pair<string, string>> get_all_arcs();
  vector<pair<string, string>> get_all_neighboring_arcs(string var)
};
