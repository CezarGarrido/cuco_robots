List<Map<String, List<Map<String, String>>>> MapByKey(String keyName,
    String newKeyName, String keyForNewName, List<Map<String, String>> input) {
  Map<String, Map<String, List<Map<String, String>>>> returnValue =
      Map<String, Map<String, List<Map<String, String>>>>();
  for (var currMap in input) {
    if (currMap.containsKey(keyName)) {
      var currKeyValue = currMap[keyName];
      var currKeyValueForNewName = currMap[keyForNewName];
      if (!returnValue.containsKey(currKeyValue)) {
        returnValue[currKeyValue] = {currKeyValue: List<Map<String, String>>()};
      }
      returnValue[currKeyValue][currKeyValue]
          .add({newKeyName: currKeyValueForNewName});
    }
  }
  return returnValue.values.toList();
}
