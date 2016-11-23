/*
** Copyright [2013-2015] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package handler
/*
func Runner(packages []string, i *WrappedParms) (string, error) {
	var outBuffer bytes.Buffer
	inputs := []string{"email=info@megam.io"} // Email is needed for urknall template's events trigger
	logWriter.Async()
	defer logWriter.Close()
	writer := io.MultiWriter(&outBuffer, os.Stdout)
	i.IfNoneAddPackages(packages)
	if h, err := NewHandler(i); err != nil {
		return "", err
	} else if err := h.Run(writer, inputs); err != nil {
		fmt.Println(err)
		return "", err
	}

	s := outBuffer.String()
	if strings.Contains(s, "failed to initiate user") {
		return s, ers.ErrUserPrivileges
	} else if strings.Contains(s, "ssh: handshake failed") {
		return s, ers.ErrAuthendication
	}
	return s, nil
}
*/
