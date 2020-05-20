package html2text

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHTML2Text(t *testing.T) {
	Convey("HTML2Text should work", t, func() {

		Convey("Links", func() {
			So(HTML2Text(`<div></div>`), ShouldEqual, "")
			So(HTML2Text(`<div>simple text</div>`), ShouldEqual, "simple text")
			So(HTML2Text(`click <a href="test">here</a> lalala`), ShouldEqual, "click here lalala")
			So(HTML2Text(`click <a class="x" href="test">here</a>`), ShouldEqual, "click here")
			So(HTML2Text(`click <a href="ents/&apos;x&apos;">here</a>`), ShouldEqual, "click here")
			So(HTML2Text(`click <a href="javascript:void(0)">here</a>`), ShouldEqual, "click here")
			So(HTML2Text(`click <a href="http://bit.ly/2n4wXRs">news</a>`), ShouldEqual, "click news")
			So(HTML2Text(
				`
				<h2><span style="color: #4b67a1;">This is a demo</span> - <span style="color: #008000;">You can edit the text! <img src="/images/smiley.png" alt="laughing" /> &hearts;</span> </h2><p>Type in the <strong>visual editor</strong> on the left or the <strong>source editor</strong> on the right and see them both change in real time.</p><p>Set up the cleaning preferences below and click the <strong style="box-shadow: 3px 3px 3px #aaa; border-radius: 5px; padding: 0 5px; background-color: #2b3a56; color: #fff;"> Clean HTML</strong> button to clean the HTML source code.</p><!--This is just a comment above the table...--><table class="demoTable" cellpadding="10"><tbody><tr style="text-align: center;"></tr><tr><td colspan="3"><strong>Convert almost any document to clean HTML in three simple steps:</strong><ol><li>Paste the content of the Document in the editor</li><li>Click the Clean HTML (optional)</li><li>Copy the generated HTML code</li></ol></td></tr></tbody></table><p><strong><span style="color: #366691; font-size: 20px; text-shadow: 4px 10px 4px #888;">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;HTML-Cleaner.com</span></strong></p><p>Please find below all the cleaning preferences, the Find And Replace tool, the Lorem-ipsum generator, the <a href="https://html-cleaner.com/case-converter/">Case Converter</a> and much more!</p><p>Don't forget to save this link into your bookmarks and share it with your friends.</p>
				`),
				ShouldEqual, "This is a demo - You can edit the text! ♥ Type in the visual editor on the left or the source editor on the right and see them both change in real time. Set up the cleaning preferences below and click the Clean HTML button to clean the HTML source code. Convert almost any document to clean HTML in three simple steps: Paste the content of the Document in the editor Click the Clean HTML (optional) Copy the generated HTML code       HTML-Cleaner.com Please find below all the cleaning preferences, the Find And Replace tool, the Lorem-ipsum generator, the Case Converter and much more! Don't forget to save this link into your bookmarks and share it with your friends.")
			So(HTML2Text(
				`
				<h2><span style="color: #4b67a1;">This is a demo</span> - <span style="color: #008000;">You can edit the text! <img src="/images/smiley.png" alt="laughing" /> &hearts;</span></h2>
				<p>Type in the <strong>visual editor</strong> on the left or the <strong>source editor</strong> on the right and see them both change in real time.</p>
				<p>Set up the cleaning preferences below and click the <strong style="box-shadow: 3px 3px 3px #aaa; border-radius: 5px; padding: 0 5px; background-color: #2b3a56; color: #fff;"> Clean HTML</strong> button to clean the HTML source code.</p>
				<!--This is just a comment above the table...-->
				<table class="demoTable" cellpadding="10">
				<tbody>
				<tr style="text-align: center;">
				<td><img src="/images/document-editors.png" alt="editors" /></td>
				<td><img src="/images/cleaning-arrow.png" alt="cleaning" /></td>
				<td><img src="/images/clean-html.png" alt="editors" width="86" height="122" /></td>
				</tr>
				<tr>
				<td colspan="3"><strong>Convert almost any document to clean HTML in three simple steps:</strong>
				<ol>
				<li>Paste the content of the Document in the editor</li>
				<li>Click the Clean HTML (optional)</li>
				<li>Copy the generated HTML code</li>
				</ol>
				</td>
				</tr>
				</tbody>
				</table>
				<p><strong><span style="color: #366691; font-size: 20px; text-shadow: 4px 10px 4px #888;">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;HTML-Cleaner.com</span></strong></p>
				<p>Please find below all the cleaning preferences, the Find And Replace tool, the Lorem-ipsum generator, the <a href="https://html-cleaner.com/case-converter/">Case Converter</a> and much more!</p>
				<p>Don't forget to save this link into your bookmarks and share it with your friends.</p>
				`),
				ShouldEqual,
				"This is a demo - You can edit the text! ♥ Type in the visual editor on the left or the source editor on the right and see them both change in real time. Set up the cleaning preferences below and click the Clean HTML button to clean the HTML source code. Convert almost any document to clean HTML in three simple steps: Paste the content of the Document in the editor Click the Clean HTML (optional) Copy the generated HTML code       HTML-Cleaner.com Please find below all the cleaning preferences, the Find And Replace tool, the Lorem-ipsum generator, the Case Converter and much more! Don't forget to save this link into your bookmarks and share it with your friends.",
			)
		})

		Convey("Inlines", func() {
			So(HTML2Text(`strong <strong>text</strong>`), ShouldEqual, "strong text")
			So(HTML2Text(`some <div id="a" class="b">div</div>`), ShouldEqual, "some div")
		})

		Convey("HTML entities", func() {
			So(HTML2Text(`two&nbsp;&nbsp;spaces`), ShouldEqual, "two  spaces")
			So(HTML2Text(`&copy; 2017 K3A`), ShouldEqual, "© 2017 K3A")
			So(HTML2Text("&lt;printtag&gt;"), ShouldEqual, "<printtag>")
			So(HTML2Text(`would you pay in &cent;, &pound;, &yen; or &euro;?`),
				ShouldEqual, "would you pay in ¢, £, ¥ or €?")
			So(HTML2Text(`Tom & Jerry is not an entity`), ShouldEqual, "Tom & Jerry is not an entity")
			So(HTML2Text(`this &neither; as you see`), ShouldEqual, "this &neither; as you see")
			So(HTML2Text(`fish &amp; chips`), ShouldEqual, "fish & chips")
			So(HTML2Text(`&quot;I'm sorry, Dave. I'm afraid I can't do that.&quot; – HAL, 2001: A Space Odyssey`), ShouldEqual, "\"I'm sorry, Dave. I'm afraid I can't do that.\" – HAL, 2001: A Space Odyssey")
			So(HTML2Text(`Google &reg;`), ShouldEqual, "Google ®")
		})

		Convey("Large Entity", func() {
			So(HTMLEntitiesToText("&abcdefghij;"), ShouldEqual, "&abcdefghij;")
		})

		Convey("Numeric HTML Entities", func() {
			So(HTMLEntitiesToText("&#39;single quotes&#39; and &#52765;"), ShouldEqual, "'single quotes' and 츝")
		})

		Convey("Full HTML structure", func() {
			So(HTML2Text(``), ShouldEqual, "")
			So(HTML2Text(`<html><head><title>Good</title></head><body>x</body>`), ShouldEqual, "x")
			So(HTML2Text(`we are not <script type="javascript"></script>interested in scripts`),
				ShouldEqual, "we are not interested in scripts")
		})
	})
}
