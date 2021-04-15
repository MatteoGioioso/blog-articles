package com.github.matteogioioso.rbijetbrainsplugin.toolWindow.runConfigurations;

import com.intellij.execution.configurations.RunConfigurationOptions;
import com.intellij.openapi.components.StoredProperty;

import java.util.HashMap;
import java.util.Map;

public class MyRunConfigurationOptions extends RunConfigurationOptions {
    // Configuration options way to store properties you can add as many properties as you need
    private final StoredProperty<Map<String, String>> environmentalVariables = map(new HashMap<String, String>())
            .provideDelegate(this, "environmentalVariables");

    private final StoredProperty<String> myCommand = string("")
            .provideDelegate(this, "myCommand");

    public Map<String, String> getEnvironmentalVariables(){
        return environmentalVariables.getValue(this);
    }

    public void setEnvironmentalVariables(Map<String, String> env){
        environmentalVariables.setValue(this, env);
    }

    public String getMyCommand() {
        return myCommand.getValue(this);
    }

    public void setMyCommand(String command) {
        myCommand.setValue(this, command);
    }
}
